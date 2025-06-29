import React, { useEffect } from 'react';
import { Layout, Menu } from 'antd';
import {
  UserOutlined,
  TeamOutlined,
  CalendarOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { Routes, Route, useNavigate, useLocation, Navigate } from 'react-router-dom';
import CoachPage from './pages/CoachPage';
import BookingPage from './pages/BookingPage';
import LoginPage from './pages/LoginPage';
import { getRole, logout } from './utils/auth';

const { Header, Sider, Content } = Layout;

const menuItems = [
  {
    key: 'coach',
    icon: <TeamOutlined />,
    label: '教练管理',
    roles: ['admin'],
  },
  {
    key: 'booking',
    icon: <CalendarOutlined />,
    label: '预约管理',
    roles: ['admin', 'coach'],
  },
  {
    key: 'logout',
    icon: <LogoutOutlined />,
    label: '退出登录',
    roles: ['admin', 'coach'],
  },
];

export default function App() {
  const navigate = useNavigate();
  const location = useLocation();
  const role = getRole();

  useEffect(() => {
    console.log('App mounted');
    console.log('Current role:', role);
    console.log('Current location:', location.pathname);
  }, []);

  console.log('App rendering, role:', role, 'pathname:', location.pathname);

  const handleMenuClick = ({ key }) => {
    console.log('Menu clicked:', key);
    if (key === 'logout') {
      logout();
      navigate('/login');
    } else {
      navigate(`/${key}`);
    }
  };

  // 如果没有角色且不在登录页，显示登录组件
  if (!role && location.pathname !== '/login') {
    return <LoginPage />;
  }

  // 如果有角色但在登录页，重定向到首页
  if (role && location.pathname === '/login') {
    return <Navigate to="/" replace />;
  }

  const filteredMenu = menuItems.filter(item => item.roles.includes(role));

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {role && (
        <Sider>
          <div style={{ color: '#fff', textAlign: 'center', padding: 16, fontWeight: 'bold' }}>
            课程预约后台
          </div>
          <Menu
            theme="dark"
            mode="inline"
            selectedKeys={[location.pathname.replace('/', '') || 'coach']}
            items={filteredMenu}
            onClick={handleMenuClick}
          />
        </Sider>
      )}
      <Layout>
        {role && (
          <Header style={{ background: '#fff', paddingLeft: 24, fontSize: 18 }}>
            欢迎，{role === 'admin' ? '管理员' : '教练'}
          </Header>
        )}
        <Content style={{ margin: 24, background: '#fff', minHeight: 360 }}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/coach" element={role ? <CoachPage /> : <Navigate to="/login" />} />
            <Route path="/booking" element={role ? <BookingPage /> : <Navigate to="/login" />} />
            <Route path="/" element={<Navigate to="/coach" replace />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
} 