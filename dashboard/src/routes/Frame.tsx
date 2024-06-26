import { Outlet } from 'react-router-dom';
import Header from '../components/Header';

export const Frame: React.FC = () => {
    return (
        <>
            <Header />
            <Outlet />
        </>
    );
};
