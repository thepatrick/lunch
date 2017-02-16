import { combineReducers } from 'redux';
import { reducer as form } from 'redux-form';
import user from './user';
import places from './places';
import locale from './locale';

const todoApp = combineReducers({
  form,
  user,
  places,
  locale,
});

export default todoApp;
