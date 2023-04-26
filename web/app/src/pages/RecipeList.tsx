import React from 'react';
import { Alert, Grid, List, ListItemButton, ListItemText, Skeleton, Stack, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { route, Routes } from '@app/router';
import { useQuery } from '@tanstack/react-query';
import { jsonApiQuery } from '@app/utils/api';

export default function RecipeList() {

  const { isLoading, isError, data } = useQuery({
    queryKey: ['recipes'],
    queryFn: jsonApiQuery('/recipes'),
  });

  const recipes = data;

  const navigate = useNavigate();
  function redirectToRecipeOptions(recipeName: string) {
    navigate(route(Routes.init, { recipeName }));
  }

  return <Grid display="flex" flexDirection="column" alignItems="center">
    <Typography
      variant="h3"
      sx= {{
        mt: 10,
        mb: 2,
      }}
    >
      Choose your recipe
    </Typography>
    {isLoading && <Stack spacing={1}>
      {/* For variant="text", adjust the height via font-size */}
      <Skeleton variant="text" sx={{ fontSize: '1rem' }} />

      {/* For other variants, adjust the size with `width` and `height` */}
      <Skeleton variant="rectangular" width={210} height={60} />
      <Skeleton variant="rectangular" width={210} height={60} />
      <Skeleton variant="rectangular" width={210} height={60} />
    </Stack>}
    {isError && <Alert severity="error">An error occurred while fetching the data</Alert>}
    {!isLoading && !isError && <List>
      {recipes.map((recipe) =>
        <ListItemButton
          key={recipe.name}
          divider={true}
          onClick={() => redirectToRecipeOptions(recipe.name)}
        >
          <ListItemText
            primary={recipe.name}
            secondary={recipe.description}
          />
        </ListItemButton>,
      )}
    </List>}
  </Grid>;
}
