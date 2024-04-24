package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoPermission error = errors.New("Sem permissão")

type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
	DeleteTx(ctx context.Context, args DeleteTxParams) error
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// inicia uma transação
// recebe uma função que pode retornar um erro
// em caso de erro, da rollback na transação
// caso não haja erro, comita a trasação
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbError := tx.Rollback()
		if rbError != nil {
			return fmt.Errorf("transaction error: %v, rolback error: %v", err, rbError)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_acount_id"`
	ToAccountID   int64 `json:"to_acount_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// executa uma transferencia entre duas contas em caso de sucesso total ou retorna um erro
func (s *SQLStore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx,
		func (q *Queries) error { // função anônima recebida pela função execTx (transfere entre duas contas)
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount: -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}
		
		// supondo que duas transações ocorram concorrentemente:
		// (1) conta1- => conta2+
		// (2) conta2- => conta1+
		// caso ocorresse: 
		// (1)conta1- seguido por (2)conta2- seguido por (1)conta2+
		// 						 					 		^
		//											 		|
		//                    					 		deadlock
		// quando ocorre a tentativa de adicionar na conta2, ela ja estava sendo bloqueada pela ação anterior.
		// verificando o id das contas, conseguimos garantir que as contas estejam "alinhadas" nas transações
		// sendo agora:
		// (1) conta1- => conta2+
		// (2) conta1+ => conta2-
		// (1)conta1- seguido por (2)conta1+ seguido por (1)conta2+ seguido por (2)conta2-
		if (args.FromAccountID < args.ToAccountID) {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
		}
		if err != nil {
			return err
		}
		
		return nil
	})

	return result, err
}

type DeleteTxParams struct {
	Owner string `json:"owner"`
	ID int64 `json:"ID"`
}

func (s *SQLStore) DeleteTx(ctx context.Context, args DeleteTxParams) error {
	return s.execTx(ctx, func(q *Queries) error {
		accountFromID, err := s.GetAccount(ctx, args.ID)
		if err != nil {
			return err
		}

		if accountFromID.Owner != args.Owner {
			return ErrNoPermission
		}

		return s.DeleteAccount(ctx, args.ID)
	})
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) { // nomear as variaveis de retorno permite ocultá-las quando retornar
	account1, err = q.AddAccoutBalance(ctx, AddAccoutBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccoutBalance(ctx, AddAccoutBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	
	return
}