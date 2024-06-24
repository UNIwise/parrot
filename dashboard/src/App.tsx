import Box from '@mui/joy/Box';
import Breadcrumbs from '@mui/joy/Breadcrumbs';
import CssBaseline from '@mui/joy/CssBaseline';
import Link from '@mui/joy/Link';
import { CssVarsProvider } from '@mui/joy/styles';
import Typography from '@mui/joy/Typography';

import ChevronRightRoundedIcon from '@mui/icons-material/ChevronRightRounded';
import HomeRoundedIcon from '@mui/icons-material/HomeRounded';

import OrderTable from './components/OrderTable';

export default function JoyOrderDashboardTemplate() {
    return (
        <CssVarsProvider disableTransitionOnChange>
            <CssBaseline />
            <Box sx={{ display: 'flex', minHeight: '100vh' }}>
                <Box
                    component="main"
                    className="MainContent"
                    sx={{
                        px: { xs: 2, md: 6 },
                        pt: {
                            xs: 'calc(12px + var(--Header-height))',
                            sm: 'calc(12px + var(--Header-height))',
                            md: 3,
                        },
                        pb: { xs: 2, sm: 2, md: 3 },
                        flex: 1,
                        display: 'flex',
                        flexDirection: 'column',
                        minWidth: 0,
                        height: '100dvh',
                        gap: 1,
                    }}
                >
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Breadcrumbs
                            size="sm"
                            aria-label="breadcrumbs"
                            separator={<ChevronRightRoundedIcon fontSize="small" />}
                            sx={{ pl: 0 }}
                        >
                            <Link underline="none" color="neutral" href="#some-link" aria-label="Home">
                                <HomeRoundedIcon />
                            </Link>

                            <Link
                                underline="hover"
                                color="neutral"
                                href="#some-link"
                                fontSize={12}
                                fontWeight={500}
                                content="Dashboard"
                            >
                                Dashboard
                            </Link>

                            <Typography color="primary" fontWeight={500} fontSize={12}>
                                Orders
                            </Typography>
                        </Breadcrumbs>
                    </Box>
                    <Box
                        sx={{
                            display: 'flex',
                            mb: 1,
                            gap: 1,
                            flexDirection: { xs: 'column', sm: 'row' },
                            alignItems: { xs: 'start', sm: 'center' },
                            flexWrap: 'wrap',
                            justifyContent: 'space-between',
                        }}
                    >
                        <Typography level="h2" component="h1">
                            Orders
                        </Typography>
                    </Box>
                    <OrderTable />
                </Box>
            </Box>
        </CssVarsProvider>
    );
}
