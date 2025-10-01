import axios, { type AxiosResponse } from 'axios';

import type { Group } from 'src/types/common.type';

export async function getGroupById(id: number): Promise<Group> {
  return axios
    .get<Group>(`${import.meta.env.VITE_API_PREFIX}/groups/${id}`)
    .then((res) => res.data);
}

export async function createGroup(name: string): Promise<{ id: number }> {
  return axios
    .post<
      { id: number },
      AxiosResponse,
      { name: string }
    >(`${import.meta.env.VITE_API_PREFIX}/groups`, { name })
    .then((res) => res.data);
}

export async function patchGroup(id: number, name: string): Promise<void> {
  await axios.patch<void, AxiosResponse, { name: string }>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${id}`,
    {
      name,
    },
  );
}
