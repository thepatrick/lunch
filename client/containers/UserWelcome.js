import { connect } from 'react-redux';
import Welcome from '../components/Welcome';

const mapStateToProps = state => (
  {
    // active: ownProps.filter === state.visibilityFilter,
    name: state.user.name,
    teamName: state.user.teamName,
    isFetching: state.user.isFetching,
    error: state.user.error,
  }
);

const UserWelcome = connect(
  mapStateToProps,
)(Welcome);

export default UserWelcome;
