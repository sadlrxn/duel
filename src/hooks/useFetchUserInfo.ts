import { useCallback } from 'react';
import { useConnection, useWallet } from '@solana/wallet-adapter-react';
import { TOKEN_PROGRAM_ID } from '@solana/spl-token';
import { useAppDispatch, useAppSelector } from 'state';
import { loadUserNfts, loadStatistics } from 'state/user/actions';
import { api, fetchUserInfo } from 'services';

const initialStatistics = {
  total_rounds: 0,
  winned_rounds: 0,
  lost_rounds: 0,
  best_streaks: 0,
  worst_streaks: 0,
  total_wagered: 0,
  max_profit: 0,
  total_profit: 0,
  private_profile: false
};

const useFetchUserInfo = () => {
  const { connection } = useConnection();
  const { publicKey } = useWallet();
  const dispatch = useAppDispatch();
  const { id } = useAppSelector(state => state.user);

  const fetchUserNfts = useCallback(async () => {
    if (!publicKey) return;

    let nftMintAddresses = [];
    try {
      const tokenAccounts = await connection.getParsedTokenAccountsByOwner(
        publicKey,
        {
          programId: TOKEN_PROGRAM_ID
        }
      );

      const nftAccounts = tokenAccounts.value.filter(
        obj =>
          obj.account.data.parsed.info.tokenAmount.uiAmount === 1 &&
          obj.account.data.parsed.info.tokenAmount.decimals === 0
      );

      nftMintAddresses = nftAccounts.map(
        element => element.account.data.parsed.info.mint
      );
    } catch (error) {
      // console.log("got error fetching NFT mint addresses!");
    }

    let undeposited = [];
    // let deposited = [];
    try {
      const { data } = await api.post('/acceptable-nfts', {
        mintAddresses: nftMintAddresses
      });

      undeposited = data;
    } catch (error) {
      // console.log("got error fetching undeposited NFTs!");
    }

    // try {
    //   const { data } = await api.post('/deposited-nfts');
    //   deposited = data;
    // } catch (error) {
    //   // console.log("got error fetching deposited NFTs!");
    // }

    dispatch(loadUserNfts({ undeposited: undeposited }));
  }, [publicKey, connection, dispatch]);

  const fetchUserStatistics = useCallback(async () => {
    if (!publicKey) return;

    try {
      const { statistics } = await fetchUserInfo({ userId: id });
      dispatch(loadStatistics(statistics));
    } catch {
      dispatch(loadStatistics(initialStatistics));
    }
  }, [publicKey, id, dispatch]);

  return {
    fetchUserNfts,
    fetchUserStatistics
  };
};

export default useFetchUserInfo;
