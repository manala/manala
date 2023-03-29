import React from 'react';
import { useParams } from 'react-router-dom';
import { Typography } from '@mui/material';

export default function ProjetConfiguration() {
  const { recipeName } = useParams();

  return <Typography>
    Recipe selected: {recipeName}
  </Typography>;
}
