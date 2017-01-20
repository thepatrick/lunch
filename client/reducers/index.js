import { combineReducers } from 'redux';
import user from './user';
import places from './places';
import locale from './locale';

const todoApp = combineReducers({
  user,
  places,
  locale,
});

export default todoApp;
