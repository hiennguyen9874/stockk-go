/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable jsx-a11y/anchor-is-valid */
import type { FC } from 'react';

import { SignInForm } from 'components/common/Forms';

const Home: FC = () => {
  return (
    <SignInForm
      onSubmit={() => {}}
      forgotPasswordLink="forgotpassword"
      signUpLink="signup"
    />
  );
};

export default Home;
