import {
  Form,
  Input,
  Checkbox,
  Link,
  Button,
  Space,
} from '@arco-design/web-react';
import { FormInstance } from '@arco-design/web-react/es/Form';
import { IconLock, IconUser } from '@arco-design/web-react/icon';
import React, { useEffect, useRef, useState } from 'react';
import axios from 'axios';
import useStorage from '@/utils/useStorage';
import useLocale from '@/utils/useLocale';
import locale from './locale';
import styles from './style/index.module.less';
import { LoginResponse, useLoginMutation } from '@/api/user';

export default function LoginForm() {
  const formRef = useRef<FormInstance>();
  const [errorMessage, setErrorMessage] = useState('');

  const [reqLogin, { isLoading }] = useLoginMutation();

  const t = useLocale(locale);

  function afterLoginSuccess(params: LoginResponse) {
    // 记录登录状态
    localStorage.setItem('userStatus', 'login');
    // 记录刷新token
    localStorage.setItem('refreshToken', params.token.refresh_token);
    // 跳转首页
    window.location.href = '/';
  }

  function onSubmitClick() {
    formRef.current.validate().then((values) => {
      reqLogin(values)
        .unwrap()
        .then((resp) => {
          afterLoginSuccess(resp.data);
        })
        .catch((e) => {
          e.data && setErrorMessage(e.data.msg);
        });
    });
  }

  return (
    <div className={styles['login-form-wrapper']}>
      <div className={styles['login-form-title']}>{t['login.form.title']}</div>
      <div className={styles['login-form-sub-title']}>
        {t['login.form.title']}
      </div>
      <div className={styles['login-form-error-msg']}>{errorMessage}</div>
      <Form
        className={styles['login-form']}
        layout="vertical"
        ref={formRef}
      >
        <Form.Item
          field="userName"
          rules={[{ required: true, message: t['login.form.userName.errMsg'] }]}
        >
          <Input
            prefix={<IconUser />}
            placeholder={t['login.form.userName.placeholder']}
            onPressEnter={onSubmitClick}
          />
        </Form.Item>
        <Form.Item
          field="password"
          rules={[{ required: true, message: t['login.form.password.errMsg'] }]}
        >
          <Input.Password
            prefix={<IconLock />}
            placeholder={t['login.form.password.placeholder']}
            onPressEnter={onSubmitClick}
          />
        </Form.Item>
        <Space size={16} direction="vertical">
          <div
            className={styles['login-form-password-actions']}
            style={{
              float: 'right',
            }}
          >
            <Link>{t['login.form.forgetPassword']}</Link>
          </div>
          <Button
            type="primary"
            long
            onClick={onSubmitClick}
            loading={isLoading}
          >
            {t['login.form.login']}
          </Button>
        </Space>
      </Form>
    </div>
  );
}
