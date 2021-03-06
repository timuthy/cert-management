/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package source

import (
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile"
	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller/reconcile/reconcilers"
)

// SlaveReconcilerType creates a slaveReconciler.
func SlaveReconcilerType(c controller.Interface) (reconcile.Interface, error) {
	reconciler := &slaveReconciler{
		controller: c,
		slaves:     c.(*reconcilers.SlaveReconciler),
	}
	return reconciler, nil
}

type slaveReconciler struct {
	reconcile.DefaultReconciler
	controller controller.Interface
	slaves     *reconcilers.SlaveReconciler
}

func (r *slaveReconciler) Start() {
	r.controller.Infof("determining dangling certificates...")
	cluster := r.controller.GetMainCluster()
	main := cluster.GetId()
	for k := range r.slaves.GetMasters(false) {
		if k.Cluster() == main {
			if _, err := cluster.GetCachedObject(k); errors.IsNotFound(err) {
				r.controller.Infof("trigger vanished origin %s", k.ObjectKey())
				r.controller.EnqueueKey(k)
			} else {
				r.controller.Debugf("found origin %s", k.ObjectKey())
			}
		}
	}
}
