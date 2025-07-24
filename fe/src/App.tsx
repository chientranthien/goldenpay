import './App.css';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Login from './Login'
import Signup from './Signup'
import Home from './Home'
import Transfer from './Transfer'
import Topup from './Topup'
import Chat from './Chat'

export default function App() {
  return (
    <Router>
      <div>
        {/*<nav>*/}
        {/*  <ul>*/}
        {/*    <li>*/}
        {/*      <Link to="/">Home</Link>*/}
        {/*    </li>*/}
        {/*    <li>*/}
        {/*      <Link to="/about">About</Link>*/}
        {/*    </li>*/}
        {/*    <li>*/}
        {/*      <Link to="/users">Users</Link>*/}
        {/*    </li>*/}
        {/*  </ul>*/}
        {/*</nav>*/}

        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/about">
            <About />
          </Route>
          <Route path="/users">
            <Users />
          </Route>
          <Route exact path="/">
            <Home />
          </Route>
          <Route path="/login">
            <Login />
          </Route>
          <Route path="/signup">
            <Signup />
          </Route>
          <Route path="/transfer">
            <Transfer />
          </Route>
          <Route path="/topup">
            <Topup />
          </Route>
          <Route path="/chat">
            <Chat />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

function About() {
  return <h2>About</h2>;
}

function Users() {
  return <h2>Users</h2>;
}
