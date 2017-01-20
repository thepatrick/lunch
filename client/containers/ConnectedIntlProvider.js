import { connect } from 'react-redux';
import { IntlProvider } from 'react-intl';

export default connect((state) => {
  const { lang, messages } = state.locale;
  return { locale: lang, key: lang, messages };
})(IntlProvider);
