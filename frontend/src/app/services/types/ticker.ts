export interface TickerResponse {
  id: number;
  symbol: string;
  exchange: string;
  full_name: string;
  short_name: string;
  type: string;
  is_active: boolean;
}

export interface TickerSnapshotResponse {
  ticker: string;
  basic_price: number;
  ceiling_price: number;
  floor_price: number;
  accumulated_vol: number;
  accumulated_val: number;
  match_price: number;
  match_qtty: number;
  highest_price: number;
  lowest_price: number;
  buy_foreign_qtty: number;
  sell_foreign_qtty: number;
  project_open: number;
  current_room: number;
  floor_code: string;
  total_room: number;
}
