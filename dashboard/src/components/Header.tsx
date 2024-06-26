import { Sheet, Typography } from '@mui/joy';
import React from 'react';
import ColorSchemeToggle from './ColorSchemeToggle';

const Header: React.FC = () => {
    return (
        <Sheet
            component="header"
            sx={{
                p: 2,
                background: `linear-gradient(
                  to right,
                  rgba(255, 0, 0, 0.7),
                  rgba(255, 165, 0, 0.7),
                  rgba(255, 255, 0, 0.7),
                  rgba(0, 128, 0, 0.7),
                  rgba(0, 0, 255, 0.7),
                  rgba(128, 0, 128, 0.7)
                )`,
                mb: 3,
                display: 'flex',
                justifyContent: 'center',
                borderRadius: '5px',
                alignItems: 'center',
                height: 64,
            }}
        >
            <Typography level="h4" sx={{ fontWeight: 'bold', color: 'primary.solidColor', fontSize: '2.2em' }}>
                Parrot
            </Typography>
            <ColorSchemeToggle sx={{ ml: 'auto' }} />
        </Sheet>
    );
};

export default Header;
