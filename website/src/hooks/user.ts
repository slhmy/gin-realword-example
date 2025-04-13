import { ModelsUser, UserApi } from "@/api/v1";
import { useEffect, useState } from "react";

export const useUser = () => {
  const [user, setUser] = useState<ModelsUser | null>(null);

  useEffect(() => {
    new UserApi().getUserMe().then((res) => {
      if (res.status === 200) {
        setUser(res.data);
      } else {
        setUser(null);
      }
    });
  }, []);

  return { user };
};
