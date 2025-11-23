import axios, { type AxiosResponse } from 'axios';

import type { Owe } from 'src/types/common.type';

export async function calculate(groupId: number): Promise<Owe[]> {
  return axios
    .post<
      Owe[],
      AxiosResponse
    >(`${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/calculate`)
    .then((res) => res.data.ious);
}
