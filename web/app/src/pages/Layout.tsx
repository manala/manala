import React from 'react';
import { AppBar, CssBaseline, Toolbar, Typography, Box, Container } from '@mui/material';
import { createTheme, ThemeProvider } from '@mui/material/styles';

const theme = createTheme();

type Props = {
  children: React.ReactNode
}

export default function Layout({ children }: Props) {
  return <ThemeProvider theme={theme}>
    <CssBaseline />
    <AppBar
      position="fixed"
      sx={{
        bgcolor: '#ef7057',
      }}
    >
      <Toolbar>
        <Typography variant="h6" color="inherit" noWrap>
          Manala Web UI
        </Typography>
      </Toolbar>
    </AppBar>
    <main>
      <Box
        sx={{
          bgcolor: 'background.paper',
          pt: 8,
          pb: 6,
        }}
      >
        <Container>
          {children}
        </Container>
      </Box>
    </main>
  </ThemeProvider>;
}
