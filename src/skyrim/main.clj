(ns skyrim.main
  (:require [clojure.core.async :as async :refer [<! <!! >! >!! chan go]]
            [clojure.data.csv :as csv]
            [clojure.java.io :as io]
            [clojure.set :as set]
            [skyrim.utils :refer [effects-left eval-potion ie-raw no-duds]]
            [clojure.core.reducers :as r]))

(defn- reduce-fn-thing [x y]
  (let [f (fn [a] (reduce + (map count (vals a))))]
    (if (> (f x) (f y))
      x
      y)))

(defn run-algo []
  (loop [ie ie-raw
         result []]
    (if (not (effects-left ie))
      result
      (do
        (let [ch (chan (count no-duds))]
          (for [n no-duds]
            (do
              (println n)
              (go
                (>! chan (eval-potion n ie)))))
          (r/reduce reduce-fn-thing ch))))))
