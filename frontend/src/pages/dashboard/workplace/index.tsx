import { Card, Typography } from '@arco-design/web-react';
import { useSelector } from 'react-redux';
import useLocale from '@/utils/useLocale';
import { selectUserInfo } from '@/store/global';
import locale from './locale';

function Workplace() {
  const t = useLocale(locale);
  const userInfo = useSelector(selectUserInfo);

  return (
    <Card>
      <Typography.Title heading={5}>
        {t['workplace.welcomeBack']}
        {userInfo.username}
      </Typography.Title>
    </Card>
  );
}

export default Workplace;
