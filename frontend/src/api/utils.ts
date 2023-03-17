// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface Response<T = any> {
  status: string;
  error: {
    code: number;
    message: string;
  } | null;
  data: T;
}
