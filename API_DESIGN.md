- login:
    support POST
      after a token will be sent

- logout:
    support POST
      the token will be ivalidanted
  
- me:
    supports GET
      gives you back selecte User infos

    supports PUT
      updates only the given inputs except for email

- registration:
    support POST
      creates new user with email and password, first name, last name, and username

- activeted/[\w]+
    supports POST
      once the user is created is created with is_active defaulted to false, and email
      will be sent with an activation link, once confirmed the account will be active
      and another final confimation email will be sent

- activeted/confirm/[\w]+
    supports POST
      once the user is created is created with is_active defaulted to false, and email
      will be sent with an activation link, once confirmed the account will be active
      and another final confimation email will be sent

- password/reset
    support POST
      given the email address an email will be sent redirecting the user to password change

- password/reset/confirm
    support POST
      given to confirm the password reset
