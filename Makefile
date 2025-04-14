GOOSE_DRIVER="postgres"
GOOSE_DBSTRING="postgresql://postgres:postgres@localhost:5432/cheapvpn"

migrate-up:
	goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

migrate-down:
	goose -dir migrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down