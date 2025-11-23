export interface Group {
  id: number;
  name: string;
}

export interface User {
  id: number;
  name: string;
  balance: number;
}

export interface Payment {
  id: number;
  description: string;
  amount: number;
  payer: User;
  payees: User[];
}

export interface Owe {
  from: number;
  to: number;
  amount: number;
}

export interface CreatePaymentRequest {
  description: string;
  amount: number;
  payer_id: number;
  payee_ids: number[];
}

export interface PatchPaymentRequest {
  id: number;
  amount: number | null;
  description: string | null;
}