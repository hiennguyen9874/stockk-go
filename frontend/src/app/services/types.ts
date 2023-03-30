export interface UserSignIn {
  email: string;
  password: string;
}

export interface UserResponse {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
  is_active: boolean;
  is_superuser: boolean;
  verified: boolean;
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  user: UserResponse;
}

export interface Response<T> {
  data: T;
  is_success: boolean;
}
