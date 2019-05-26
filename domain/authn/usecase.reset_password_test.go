package authn_test

//
// func TestResetPassword(t *testing.T) {
//         repo := authn.NewRepository(db)
//         tokenTTL := 1 * time.Minute
//
//         {
//                 useCase := authn.RegisterUseCase(repo)
//         }
//         {
//                 useCase := authn.NewRecoverPasswordUseCase(repo, tokenTTL)
//                 useCase.RecoverPassword(context.TODO(), authn.RecoverPasswordUseCase{
//                         Email: "john.doe@mail.com",
//                 })
//         }
//
//         resetPasswordUseCase := authn.NewResetPasswordUseCase(repo, tokenTTL)
//
//         req := authn.ResetPasswordRequest{
//                 Token:           "xyz",
//                 Password:        "1",
//                 ConfirmPassword: "2",
//         }
//         resetPasswordUseCase.ResetPassword(context.TODO(), req)
// }
