export interface StockSnapshotResponse {
  ticker: string;
  ref_price: number;
  ceil_price: number;
  floor_price: number;
  tlt_vol: number;
  tlt_val: number;
  price_b3: number;
  price_b2: number;
  price_b1: number;
  vol_b3: number;
  vol_b2: number;
  vol_b1: number;
  price: number;
  vol: number;
  price_s3: number;
  price_s2: number;
  price_s1: number;
  vol_s3: number;
  vol_s2: number;
  vol_s1: number;
  high: number;
  low: number;
  buy_foreign: number;
  sell_foreign: number;
}
