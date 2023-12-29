<p>Packages:</p>
<ul>
<li>
<a href="#lifecycle.ironcore.dev%2fv1alpha1">lifecycle.ironcore.dev/v1alpha1</a>
</li>
</ul>
<h2 id="lifecycle.ironcore.dev/v1alpha1">lifecycle.ironcore.dev/v1alpha1</h2>
<div>
<p>Package v1alpha1 is the v1alpha1 version of the API.</p>
</div>
Resource Types:
<ul></ul>
<h3 id="lifecycle.ironcore.dev/v1alpha1.AvailablePackageVersions">AvailablePackageVersions
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineTypeStatus">MachineTypeStatus</a>)
</p>
<div>
<p>AvailablePackageVersions defines a number of versions for concrete firmware package.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name reflects the name of the firmware package</p>
</td>
</tr>
<tr>
<td>
<code>versions</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Versions reflects the list of discovered package versions available for installation.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.Machine">Machine
</h3>
<div>
<p>Machine is the Schema for the machines API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineSpec">
MachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>machineTypeRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>MachineTypeRef contain reference to MachineType object.</p>
</td>
</tr>
<tr>
<td>
<code>oobMachineRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>OOBMachineRef contains reference to OOB machine object.</p>
</td>
</tr>
<tr>
<td>
<code>scanPeriod</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>ScanPeriod defines the interval between scans.</p>
</td>
</tr>
<tr>
<td>
<code>packages</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.PackageVersion">
[]PackageVersion
</a>
</em>
</td>
<td>
<p>Packages defines the list of package versions to install.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineStatus">
MachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineGroup">MachineGroup
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineTypeSpec">MachineTypeSpec</a>)
</p>
<div>
<p>MachineGroup defines group of Machine objects filtered by label selector
and a list of firmware packages versions which should be installed by default.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>machineSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>MachineSelector defines native kubernetes label selector to apply to Machine objects.</p>
</td>
</tr>
<tr>
<td>
<code>packages</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.PackageVersion">
[]PackageVersion
</a>
</em>
</td>
<td>
<p>Packages defines default firmware package versions for the group of Machine objects.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineSpec">MachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.Machine">Machine</a>)
</p>
<div>
<p>MachineSpec defines the desired state of Machine.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>machineTypeRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>MachineTypeRef contain reference to MachineType object.</p>
</td>
</tr>
<tr>
<td>
<code>oobMachineRef</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core">
Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>OOBMachineRef contains reference to OOB machine object.</p>
</td>
</tr>
<tr>
<td>
<code>scanPeriod</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>ScanPeriod defines the interval between scans.</p>
</td>
</tr>
<tr>
<td>
<code>packages</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.PackageVersion">
[]PackageVersion
</a>
</em>
</td>
<td>
<p>Packages defines the list of package versions to install.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineStatus">MachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.Machine">Machine</a>)
</p>
<div>
<p>MachineStatus defines the observed state of Machine.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>lastScanTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastScanTime reflects the timestamp when the last scan job for installed
firmware versions was performed.</p>
</td>
</tr>
<tr>
<td>
<code>lastScanResult</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.ScanResult">
ScanResult
</a>
</em>
</td>
<td>
<p>LastScanResult reflects either success or failure of the last scan job.</p>
</td>
</tr>
<tr>
<td>
<code>installedPackages</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.PackageVersion">
[]PackageVersion
</a>
</em>
</td>
<td>
<p>InstalledPackages reflects the versions of installed firmware packages.</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message contains verbose message explaining current state</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineType">MachineType
</h3>
<div>
<p>MachineType is the Schema for the machinetypes API.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineTypeSpec">
MachineTypeSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>manufacturer</code><br/>
<em>
string
</em>
</td>
<td>
<p>Manufacturer refers to manufacturer, e.g. Lenovo, Dell etc.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>Type refers to machine type, e.g. 7z21 for Lenovo, R440 for Dell etc.</p>
</td>
</tr>
<tr>
<td>
<code>scanPeriod</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>ScanPeriod defines the interval between scans.</p>
</td>
</tr>
<tr>
<td>
<code>machineGroups</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineGroup">
[]MachineGroup
</a>
</em>
</td>
<td>
<p>MachineGroups defines list of MachineGroup</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineTypeStatus">
MachineTypeStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineTypeSpec">MachineTypeSpec
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineType">MachineType</a>)
</p>
<div>
<p>MachineTypeSpec defines the desired state of MachineType.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>manufacturer</code><br/>
<em>
string
</em>
</td>
<td>
<p>Manufacturer refers to manufacturer, e.g. Lenovo, Dell etc.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
<p>Type refers to machine type, e.g. 7z21 for Lenovo, R440 for Dell etc.</p>
</td>
</tr>
<tr>
<td>
<code>scanPeriod</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#duration-v1-meta">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>ScanPeriod defines the interval between scans.</p>
</td>
</tr>
<tr>
<td>
<code>machineGroups</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.MachineGroup">
[]MachineGroup
</a>
</em>
</td>
<td>
<p>MachineGroups defines list of MachineGroup</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.MachineTypeStatus">MachineTypeStatus
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineType">MachineType</a>)
</p>
<div>
<p>MachineTypeStatus defines the observed state of MachineType.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>lastScanTime</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>LastScanTime reflects the timestamp when the last scan of available packages was done.</p>
</td>
</tr>
<tr>
<td>
<code>lastScanResult</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.ScanResult">
ScanResult
</a>
</em>
</td>
<td>
<p>LastScanResult reflects the result of the last scan.</p>
</td>
</tr>
<tr>
<td>
<code>availablePackages</code><br/>
<em>
<a href="#lifecycle.ironcore.dev/v1alpha1.AvailablePackageVersions">
[]AvailablePackageVersions
</a>
</em>
</td>
<td>
<p>AvailablePackages reflects the list of AvailablePackageVersion</p>
</td>
</tr>
<tr>
<td>
<code>message</code><br/>
<em>
string
</em>
</td>
<td>
<p>Message contains verbose message explaining current state</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.PackageVersion">PackageVersion
</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineGroup">MachineGroup</a>, <a href="#lifecycle.ironcore.dev/v1alpha1.MachineSpec">MachineSpec</a>, <a href="#lifecycle.ironcore.dev/v1alpha1.MachineStatus">MachineStatus</a>)
</p>
<div>
<p>PackageVersion defines the concrete package version item.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name defines the name of the firmware package.</p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<p>Version defines the version of the firmware package.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="lifecycle.ironcore.dev/v1alpha1.ScanResult">ScanResult
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#lifecycle.ironcore.dev/v1alpha1.MachineStatus">MachineStatus</a>, <a href="#lifecycle.ironcore.dev/v1alpha1.MachineTypeStatus">MachineTypeStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Failure&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;Success&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
</em></p>
