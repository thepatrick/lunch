import { connect } from 'react-redux';
import App from '../components/App';

const mapStateToProps = state => (
  {
    userFetching: state.user.isFetching,
    userError: state.user.error,
  }
);

const AppContainer = connect(
  mapStateToProps,
)(App);

export default AppContainer;
