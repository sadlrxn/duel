import { AxiosError } from "axios";
import useSWR from "swr";
import { fetchUserInfo, FetchUserInfoResponse } from "services";

export default function useUserInfo({
  userId,
  userName,
}: {
  userId?: number;
  userName?: string;
}) {
  return useSWR<FetchUserInfoResponse, AxiosError>(
    `User Info: id:${userId} name:${userName}`,
    () => fetchUserInfo({ userId, userName })
  );
}
