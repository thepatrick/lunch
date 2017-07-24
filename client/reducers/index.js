import { combineReducers } from 'redux';
import { reducer as form } from 'redux-form';

import { routerReducer } from 'react-router-redux';

import user from './user';
import places from './places';
import locale from './locale';

const todoApp = combineReducers({
  form,
  user,
  places,
  locale,
  routing: routerReducer,
});

export default todoApp;
