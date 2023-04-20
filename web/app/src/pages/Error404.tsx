import React from 'react';
import Layout from '@app/pages/Layout';
import { Typography, Button, Grid } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { route, Routes } from '@app/router';

export default function Error404() {
  const navigate = useNavigate();

  function redirectToHome() {
    navigate(route(Routes.recipes));
  }

  return <Layout>
    <Grid display="flex" flexDirection="column" alignItems="center">
      <Typography
        variant="h3"
        sx= {{
          mt: 10,
          mb: 2,
        }}
      >
        Error 404
      </Typography>
      <Typography variant="body1">This page doesn&apos;t exist.</Typography>
      <Button
        onClick={() => redirectToHome()}
        variant="contained"
        sx={{
          mt: 5,
          bgcolor: '#ef7057',
          '&:hover': {
            backgroundColor: '#e86a52',
          },
        }}
      >
        Back to homepage
      </Button>
    </Grid>
  </Layout>;
}
