import React, { useState } from 'react';
import MobileCoachPage from './MobileCoachPage';
import MobileBookingPage from './MobileBookingPage';
import MobileCoachProfilePage from './MobileCoachProfilePage';
import { Button } from 'antd';
import { logout, getRole } from '../utils/auth';
import './MobileAdminPage.css';

const TABS = [
  { key: 'coach', label: '教练管理' },
  { key: 'booking', label: '预约管理' },
  { key: 'profile', label: '个人中心' },
];

export default function MobileAdminPage() {
  const [tab, setTab] = useState('coach');
  const role = getRole();

  // 退出登录
  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  return (
    <div className="mobile-admin-root">
      {/* 顶部导航 */}
      <div className="mobile-admin-navbar">
        <span className="mobile-admin-navbar-title">{TABS.find(t => t.key === tab)?.label}</span>
        <Button type="link" className="mobile-admin-logout-btn" onClick={handleLogout}>退出登录</Button>
      </div>
      {/* 功能切换按钮 */}
      <div className="mobile-admin-tabs">
        {TABS.map(t => (
          <Button
            key={t.key}
            type={tab === t.key ? 'primary' : 'default'}
            className="mobile-admin-tab-btn"
            onClick={() => setTab(t.key)}
            style={{ marginRight: 12 }}
          >
            {t.label}
          </Button>
        ))}
      </div>
      {/* 子页面内容 */}
      <div className="mobile-admin-content">
        {tab === 'coach' && <MobileCoachPage />}
        {tab === 'booking' && <MobileBookingPage />}
        {tab === 'profile' && <MobileCoachProfilePage />}
      </div>
    </div>
  );
} 