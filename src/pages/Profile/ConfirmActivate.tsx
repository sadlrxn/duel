import { useModal } from 'components';
import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { api } from 'services';
import toast from 'utils/toast';
import { Link, useNavigate } from 'react-router-dom';
import ConfirmActivateModal from './components/ConfirmActivateModal';

export default function ConfirmActivate() {
  const { referralCode } = useParams();
  const navigate = useNavigate();

  const activateReferralCode = async () => {
    try {
      await api.post('/affiliate/activate', { code: referralCode });
      toast.success('Referral code updated!');
    } catch (err: any) {
      if (err.response.status === 406) {
        if (err.response.data.errorCode === 13031)
          toast.error('You can’t activate your own code.');
        else if (err.response.data.errorCode === 13032)
          toast.error('This referral code doesn’t exist.');
        else if (err.response.data.errorCode === 13034)
          toast.error('You can’t activate code after 24 hours from sign up');
        else if (err.response.data.errorCode === 13033)
          toast.error('Failed! There was an error activating referral code.');
      } else if (err.response.status === 503) {
        toast.error('This function is blocked by admin.');
      } else
        toast.error('Failed! There was an error activating referral code.');
    } finally {
      navigate('/profile?tab=referral');
    }
  };

  const [onConfirmActivate] = useModal(
    <ConfirmActivateModal
      code={referralCode ? referralCode : ''}
      onActivate={activateReferralCode}
    />
  );

  useEffect(() => {
    if (!referralCode) return;

    onConfirmActivate();
  }, [referralCode]);

  return <div></div>;
}
