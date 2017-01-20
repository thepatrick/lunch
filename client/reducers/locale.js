import en from 'react-intl/locale-data/en';
import { addLocaleData } from 'react-intl';

addLocaleData(en);

const ENGLISH_STATE = {
  lang: 'en',
  messages: {},
};

const initialState = {
  lang: ENGLISH_STATE.lang,
  messages: ENGLISH_STATE.messages,
};

const locale = (state = initialState, action) => {
  switch (action.type) {
    default:
      return state;
  }
};

export default locale;
