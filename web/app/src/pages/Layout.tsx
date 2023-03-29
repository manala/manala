import React from 'react';
import { Outlet } from 'react-router-dom';
import { AppBar, CssBaseline, Toolbar, Typography } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';

const theme = createTheme();

export default function Layout() {
  return <ThemeProvider theme={theme}>
    <CssBaseline />
    <AppBar position="fixed" style={{ background: '#EF7057' }}>
      <Toolbar>
        <Typography variant="h6" color="inherit" noWrap>
          Manala Web UI
        </Typography>
      </Toolbar>
    </AppBar>
    <main>
      <Outlet />
    </main>
  </ThemeProvider>;
}
