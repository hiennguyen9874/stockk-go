import type { FC } from 'react';
import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';

import { useLoginMutation } from 'app/services/auth';
import { SignInForm } from 'components/common/Forms';

interface LoginProps {
  homeLink: string;
  forgotPasswordLink: string;
  signUpLink: string;
}

const Login: FC<LoginProps> = ({
  homeLink,
  forgotPasswordLink,
  signUpLink,
}) => {
  const navigate = useNavigate();

  const [login] = useLoginMutation();

  const handleLogin = useCallback(
    async (email: string, password: string) => {
      try {
        await login({ email, password }).unwrap();
        // Being that the result is handled in extraReducers in authSlice,
        // we know that we're authenticated after this, so the user
        // and token will be present in the store

        navigate(homeLink);
      } catch (err) {
        console.log(err);
      }
    },
    [homeLink, login, navigate]
  );

  return (
    <SignInForm
      onSubmit={handleLogin}
      forgotPasswordLink={forgotPasswordLink}
      signUpLink={signUpLink}
    />
  );
};

export default Login;
