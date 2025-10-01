import axios, { type AxiosResponse } from 'axios';

import type {
  CreatePaymentRequest,
  PatchPaymentRequest,
  Payment,
} from 'src/types/common.type';

export async function getPaymentsByGroupId(
  groupId: number,
): Promise<Payment[]> {
  return axios
    .get<
      Payment[]
    >(`${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/payments`)
    .then((res) => res.data);
}

export async function addPaymentToGroup(
  id: number,
  data: CreatePaymentRequest,
): Promise<{ id: number }> {
  return axios
    .post<
      { id: number },
      AxiosResponse,
      CreatePaymentRequest
    >(`${import.meta.env.VITE_API_PREFIX}/groups/${id}/payments`, data)
    .then((res) => res.data);
}

export async function patchPayment(
  groupId: number,
  id: number,
  data: PatchPaymentRequest,
): Promise<void> {
  await axios.patch<void, AxiosResponse, PatchPaymentRequest>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/payments/${id}`,
    data,
  );
}

export async function deletePayment(
  groupId: number,
  id: number,
): Promise<void> {
  await axios.delete<void, AxiosResponse>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/payments/${id}`,
  );
}

export async function deleteAllPayments(groupId: number): Promise<void> {
  await axios.delete<void, AxiosResponse>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/payments`,
  );
}
