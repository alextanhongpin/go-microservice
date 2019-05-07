package service

// Package service contains all the service (a.k.a usecases) for the
// application. We use the convention <action>_<resource> or
// <actor>_<action>_<resource>svc if actor is present.

// Grouping services by resources is easy to begin, but hard to manage later as the service grow due to lack of context.
// Say we have a dashboard application where the Ops (Operations a.k.a Admin)
// can add a new user, invite them, change the user's role etc. Ops can also
// manage promotion banners for display in the webui that the user can view.
// The initial idea may be to group the service by resources:

// - opssvc: CreateUser, ViewUsers, ViewUser, DeleteUser, InviteUser, OpsLogin,
// OpsRegister
// - usersvc: UpdateInfo, OpsLogin, OpsRegister
// - adsvc: CreateAds, UpdateAds, OpsViewAds, UserViewAds

// But we are losing a lot of context here, especially when there are more
// services for ops to manage. Note that the adsvc is now handling multiple
// roles too - one for user view and another for ops. Since both user and ops
// have different permissions, they can definitely see different info. The
// usersvc is a little ambiguous too, since it raises the question - can one
// user update another user? A better way to group the services is:

// - ops_login: Login, Register
// - ops_manage_user: CRUD.
// - ops_manage_userrole: CRUD for roles.
// - ops_manage_ads: It is a CRUD service, but we can add more context to it such
// as CreatePersonalizedAd, CreatePublicAd etc. Be specific.
// - user_login: Login, Register implementation is different than that of ops.
// - user_manage_self: UserInfo, UpdateInfo
// - user_view_ads: ViewAds, ClickAds.

// There are two caveats here:
// - go will warn if there's underscore in package name. But we are designing
// an application, not package...
// - we have more services now - but they are specific, independent by one
// another and grouped by features. It will definitely be more maintainable as
// the application grow. To remove a feature, one can easily just delete the
// folder too!
