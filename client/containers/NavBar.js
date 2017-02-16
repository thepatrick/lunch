import { connect } from 'react-redux';
import Nav from '../components/Nav';

const mapStateToProps = state => (
  {
    // active: ownProps.filter === state.visibilityFilter,
    name: state.user.name,
    teamName: state.user.teamName,
    isFetching: state.user.isFetching,
    error: state.user.error,
  }
);

const NavBar = connect(
  mapStateToProps,
)(Nav);

export default NavBar;
