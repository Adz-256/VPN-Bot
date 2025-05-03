package api

//func (a *API) StartExpireCRON() {
//	t := time.NewTicker(24 * time.Hour)
//	defer t.Stop()
//	for {
//		<-t.C
//		accs, err := a.sub.GetExpiredAccounts(context.Background())
//		if err != nil {
//			slog.Error("GetExpiredAccounts error", slog.Any("error", err))
//		}
//
//		for _, acc := range *accs {
//			err = a.sub.Block(context.Background(), acc.PublicKey)
//			if err != nil {
//				slog.Error("Block error", slog.Any("error", err))
//			}
//		}
//	}
//}
