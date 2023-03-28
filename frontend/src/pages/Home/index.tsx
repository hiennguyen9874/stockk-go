/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable jsx-a11y/anchor-is-valid */
import type { FC } from 'react';

import Login from 'features/auth/Login';

const Home: FC = () => {
  return (
    <Login
      homeLink="/"
      forgotPasswordLink="/forgotpassword"
      signUpLink="/signup"
    />
  );
};

export default Home;
