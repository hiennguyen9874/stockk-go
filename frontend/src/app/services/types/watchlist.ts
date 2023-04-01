export interface WatchListCreate {
  name: string;
}

export interface WatchListResponse {
  id: number;
  created_at: string;
  updated_at: string;
  name: string;
  tickers: string[];
  owner_id: number;
}

export interface WatchListUpdate {
  id: number;
  name: string;
  tickers: string[];
}
