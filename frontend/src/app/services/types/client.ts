export interface ClientCreate {
  current_ticker: string;
  current_resolution: string;
}

export interface ClientResponse {
  id: number;
  created_at: string;
  updated_at: string;
  current_ticker: string;
  current_resolution: string;
  owner_id: number;
}

export interface ClientUpdate {
  id: number;
  current_ticker: string;
  current_resolution: string;
}
