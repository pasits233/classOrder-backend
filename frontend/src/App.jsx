import React from 'react';
import { Layout, Menu } from 'antd';
import {
  UserOutlined,
  TeamOutlined,
  CalendarOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
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

  if (!role && location.pathname !== '/login') {
    navigate('/login');
    return null;
  }

  if (role && location.pathname === '/login') {
    navigate('/');
    return null;
  }

  const filteredMenu = menuItems.filter(item => item.roles.includes(role));

  const handleMenuClick = ({ key }) => {
    if (key === 'logout') {
      logout();
      navigate('/login');
    } else {
      navigate(`/${key}`);
    }
  };

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
            <Route path="/coach" element={<CoachPage />} />
            <Route path="/booking" element={<BookingPage />} />
            <Route path="*" element={<CoachPage />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
} 