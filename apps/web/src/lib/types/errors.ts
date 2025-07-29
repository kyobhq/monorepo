export interface APIError {
  status: number;
  code: string;
  cause: string;
  message: string;
}

export interface APIDefaultError {
  code: string;
  error?: string;
  cause?: any;
}
