import React, { useEffect, useState } from "react";
import UserService from "./api/userService";
import Cookies from "js-cookie";
import { useHistory } from "react-router-dom";

interface CommonService {
  SetFormData: (e: React.ChangeEvent<HTMLInputElement>, setFormData: React.Dispatch<React.SetStateAction<any>>) => void;
}

const common: CommonService = {
  SetFormData: function (e: React.ChangeEvent<HTMLInputElement>, setFormData: React.Dispatch<React.SetStateAction<any>>): void {
    setFormData((prev: any) => {
      return {
        ...prev,
        [e.target.name]: e.target.value
      }
    })
  }
}

export default function Common(): CommonService {
  return common
}

interface Code {
  id: number;
  msg: string;
}

export function useRedirectToLoginIfNotAuthenticated(): { isAuthenticated: boolean | null } {
  const history = useHistory()
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null)

  useEffect(() => {
    const token = Cookies.get(Constants.CookieToken);

    if (token === "" || token === undefined || token === null) {
      history.push('/login');
      setIsAuthenticated(false);
      return;
    }

    UserService().Authz().then((c: Code) => {
      if (c.id !== 0) {
        history.push('/login');
        setIsAuthenticated(false);
      } else {
        setIsAuthenticated(true);
      }
    }).catch((error) => {
      history.push('/login');
      setIsAuthenticated(false);
    });
  }, [history]);

  return { isAuthenticated };
}

export function useRedirectToHomeIfAlreadyAuthenticated(): { isCheckingAuth: boolean } {
  const history = useHistory()
  const [isCheckingAuth, setIsCheckingAuth] = useState<boolean>(true);

  useEffect(() => {
    console.log("🏠 Checking if already authenticated for redirect to home...")

    const token = Cookies.get(Constants.CookieToken);
    console.log("🍪 Token from cookie:", token)

    if (token && token !== "") {
      console.log("🔐 Token found, calling authz API...")
      UserService().Authz().then((c: Code) => {
        console.log("🔐 Authz API response code:", c)
        if (c.id === 0) {
          console.log("✅ Already authenticated, redirecting to home")
          history.push('/');
        } else {
          console.log("❌ Invalid token, removing cookie")
          Cookies.remove(Constants.CookieToken)
        }
        setIsCheckingAuth(false);
      }).catch((error) => {
        console.log("❌ Authz API error:", error)
        Cookies.remove(Constants.CookieToken)
        setIsCheckingAuth(false);
      });
    } else {
      console.log("ℹ️ No token found, staying on current page")
      setIsCheckingAuth(false);
    }
  }, [history]);

  return { isCheckingAuth };
}

export function RemoveTokenCookie(): void {
  Cookies.remove(Constants.CookieToken)
}

export function GetUserIdFromCookie(): number | undefined {
  return parseInt(Cookies.get(Constants.CookieUserId) || "0")
}

export const Constants = {
  CookieToken: "token",
  CookieUserId: "uid"
}

export const CSS = {
  FormCol: "col-lg-6 col-md-8 col-sm-12 col-xs-12"
}