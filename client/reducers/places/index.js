import { combineReducers } from 'redux';

import placesById from './placesById';
import placeList from './placeList';
import isFetching from './isFetching';
import error from './error';

const places = combineReducers({
  isFetching,
  error,
  placesById,
  placeList,
});

export default places;
