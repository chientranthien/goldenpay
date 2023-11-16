import {useCallback, useEffect} from "react";
import UserService from "./api/userService";
import Cookies from "js-cookie";
import {useHistory} from "react-router-dom";

const common = {}


common.SetFormData = function (e, setFormData) {
  setFormData(prev => {
    return {
      ...prev,
      [e.target.name]: e.target.value
    }
  })
}

export default function Common() {
  return common
}

export function RedirectToHomeIfAlreadyAuthenticated() {
  const history = useHistory()
  const memoizedAuthz = useCallback(() => {
    const token = Cookies.get(Constants.CookieToken);
    if (token !== "") {
      UserService().Authz().then((c) => {
        if (c.id == 0) {
          history.push('/');
          return
        }

        Cookies.remove(Constants.CookieToken)
      });
    }
  }, []);

  useEffect(() => {
    memoizedAuthz();
  }, [])
}

export const Constants = {
  CookieToken: "token"
}

export const CSS = {
  FormCol: "col-lg-6 col-md-8 col-sm-12 col-xs-12"
}