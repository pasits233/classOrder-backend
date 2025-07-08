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
import MobileLoginPage from './pages/MobileLoginPage';
import MobileCoachPage from './pages/MobileCoachPage';
import MobileBookingPage from './pages/MobileBookingPage';
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

function useIsMobile() {
  const [isMobile, setIsMobile] = React.useState(window.innerWidth < 600);
  React.useEffect(() => {
    const onResize = () => setIsMobile(window.innerWidth < 600);
    window.addEventListener('resize', onResize);
    return () => window.removeEventListener('resize', onResize);
  }, []);
  return isMobile;
}

export default function App() {
  const navigate = useNavigate();
  const location = useLocation();
  const role = getRole();
  const isMobile = useIsMobile();

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
    return isMobile ? <MobileLoginPage /> : <LoginPage />;
  }

  // 如果有角色但在登录页，重定向到首页
  if (role && location.pathname === '/login') {
    // 管理员跳转到/coach，教练跳转到/booking
    return <Navigate to={role === 'admin' ? '/coach' : '/booking'} replace />;
  }

  // 教练禁止访问教练管理页面，强制跳转到预约管理
  if (role === 'coach' && location.pathname !== '/booking') {
    return <Navigate to="/booking" replace />;
  }

  if (isMobile) {
    // 移动端：只渲染移动端页面，无PC端布局
    return (
      <Routes>
        <Route path="/login" element={<MobileLoginPage />} />
        <Route path="/coach" element={role === 'admin' ? <MobileCoachPage /> : <Navigate to="/booking" />} />
        <Route path="/booking" element={<MobileBookingPage />} />
        <Route path="/" element={<Navigate to={role === 'admin' ? '/coach' : '/booking'} replace />} />
        <Route path="*" element={<Navigate to={role === 'admin' ? '/coach' : '/booking'} replace />} />
      </Routes>
    );
  }

  // PC端：原有Layout结构
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
            <Route path="/coach" element={role === 'admin' ? <CoachPage /> : <Navigate to="/booking" />} />
            <Route path="/booking" element={<BookingPage />} />
            <Route path="/" element={<Navigate to={role === 'admin' ? '/coach' : '/booking'} replace />} />
            <Route path="*" element={<Navigate to={role === 'admin' ? '/coach' : '/booking'} replace />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  );
} 