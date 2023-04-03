export interface TickerResponse {
  id: number;
  symbol: string;
  exchange: string;
  full_name: string;
  short_name: string;
  type: string;
  is_active: boolean;
}
