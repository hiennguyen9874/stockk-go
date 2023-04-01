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
